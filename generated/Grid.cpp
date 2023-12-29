#include "Grid.h"
#include <math.h>

namespace glow
{
#ifndef MICRO_CONTROLLER

  std::string Grid::make_code()
  {
    std::stringstream s;
    s << "{" << length << ","
      << rows << ","
      << origin << ","
      << orientation << "}";
    return s.str();
  }

  std::string Grid::keys[Grid::KEY_COUNT] = {
      "length",
      "rows",
      "origin",
      "orientation",
  };

  std::string Grid::origin_keys[ORIGIN_COUNT] = {
      "top left",
      "top right",
      "bottom left",
      "bottom right",
  };

  std::unordered_map<std::string, uint16_t> Grid::origin_map = {
      {origin_keys[TopLeft], TopLeft},
      {origin_keys[TopRight], TopRight},
      {origin_keys[BottomLeft], BottomLeft},
      {origin_keys[BottomRight], BottomRight},
  };

  std::string Grid::orientation_keys[ORIENTATION_COUNT] = {
      "horizontal",
      "vertical",
      "diagonal",
      "centred",
  };

  std::unordered_map<std::string, uint16_t> Grid::orientation_map = {
      {orientation_keys[Horizontal], Horizontal},
      {orientation_keys[Vertical], Vertical},
      {orientation_keys[Diagonal], Diagonal},
      {orientation_keys[Centred], Centred},
  };

#endif

  uint16_t Grid::adjust_bounds(float bound)
  {
    uint16_t scaled = static_cast<uint16_t>(round(bound));
    if (orientation == Horizontal)
    {
      return (scaled / columns) * columns;
    }

    return (scaled / rows) * rows;
  }

  bool Grid::setup()
  {
    if (length == 0)
    {
      return false;
    }

    if (rows == 0)
    {
      rows = 1;
    }

    columns = length / rows;

    if (orientation == Centred)
    {
      setup_centred(rows, columns);
    }
    else
    {
      setup_diagonal(rows, columns);
    }

    return true;
  }

  bool Grid::setup(uint16_t p_length, uint16_t p_rows, uint8_t p_origin, uint8_t p_orientation)
  {
    length = p_length;
    rows = p_rows;
    origin = p_origin;
    orientation = p_orientation;
    return setup();
  }

  uint16_t Grid::map(uint16_t index)
  {
    uint16_t offset = index;
    if (orientation == Diagonal)
    {
      offset = map_diagonal(index);
    }
    else if (orientation == Vertical)
    {
      offset = map_columns(index);
    }
    else if (orientation == Centred)
    {
      offset = map_centred(index);
    }
    return map_to_origin(offset);
  }

  uint16_t Grid::map_diagonal_top(uint16_t index)
  {
    uint16_t offset = 0;
    uint16_t start = 0;
    while (start < index)
    {
      offset++;
      start += offset;
    }

    if (start == index)
    {
      return offset;
    }

    start -= offset;
    offset--;
    return offset + (index - start) * (columns - 1);
  }

  uint16_t Grid::map_diagonal_bottom(uint16_t index)
  {
    uint16_t offset = 2 * columns - 1;
    uint16_t start = last_edge + 1;
    uint16_t increment = rows;
    while (start < index)
    {
      offset += columns;
      increment--;
      start += increment;
    }

    if (start == index)
    {
      return offset;
    }

    start -= increment;
    offset -= columns;
    return offset + (index - start) * (columns - 1);
  }

  uint16_t Grid::map_diagonal(uint16_t index)
  {
    if (columns < 3)
    {
      return index;
    }

    if (index < first_edge)
    {
      return map_diagonal_top(index);
    }

    if (index <= last_edge)
    {
      return map_diagonal_middle(index);
    }

    return map_diagonal_bottom(index);
  }

  uint16_t Grid::map_to_origin(uint16_t offset)
  {
    if (origin == BottomRight)
    {
      return length - offset - 1;
    }

    div_t point = div(offset, columns);

    // point.quot = row
    // point.rem = col

    if (origin == BottomLeft)
    {
      return (rows - point.quot - 1) * columns +
             point.rem;
    }

    if (origin == TopRight)
    {
      return point.quot * columns +
             (columns - point.rem - 1);
    }
    return offset;
  }

  void Grid::setup_diagonal(uint16_t rows, uint16_t columns)
  {
    first_edge = 0;
    uint16_t lesser = std::min(rows, columns);
    for (uint16_t i = 0; i < lesser; i++)
    {
      first_edge += i;
    }

    centre = lesser - 1;
    last_edge = first_edge +
                (columns - lesser) * rows +
                rows - 1;
  }

  void Grid::setup_centred(uint16_t rows, uint16_t columns)
  {
    centre = ((rows - 1) >> 1) * columns + ((columns - 1) >> 1);
    if (rows == columns)
    {
      ring_status = PIVOT_SQUARE;
      ring_count = rows >> 1;
      first_edge =
          last_edge =
              first_offset = length;
      return;
    }

    if (rows < columns)
    {
      ring_status = (rows & 1) ? PIVOT_COLUMNS
                               : PIVOT_COLUMNS | PIVOT_UNEVEN;
      ring_count = (rows - 1) >> 1;
    }
    else if (rows > columns)
    {
      ring_status = (columns & 1) ? PIVOT_ROWS
                                  : PIVOT_ROWS | PIVOT_UNEVEN;
      ring_count = (columns - 1) >> 1;
    }

    last_edge = length;
    first_edge = 1;
    uint16_t ring_length = 8;
    for (uint16_t i = 0; i < ring_count; i++)
    {
      first_edge += ring_length;
      ring_length += 8;
    }

    first_offset = ((columns - 1) >> 1) + ring_count + 1;
    last_offset = (rows - 1) * columns +
                  ((columns - 1) >> 1) - ring_count - 1;
  }

  uint16_t Grid::map_centred_edge(uint16_t index)
  {
    index -= first_edge;
    uint16_t offset = first_offset;
    if (index < rows)
    {
      return offset + index * columns;
    }

    offset += (rows - 1) * columns;
    index -= rows - 1;

    if (ring_status & PIVOT_UNEVEN)
    {
      uint16_t side_width = (ring_count + 1) << 1;
      if (index < side_width)
      {
        return offset - index;
      }
      index -= side_width;
    }

    if (index < rows)
    {
      offset = last_offset - index * columns;
      return offset;
    }

    index -= rows;
    uint16_t side = index / rows;
    index -= side * rows;

    if (side & 1)
    {
      side = side / 2 + 1;
      offset = last_offset - side - columns * index;
      return offset;
    }

    side = side / 2 + 1;
    offset = first_offset + side + columns * index;
    return offset;
  }

  uint16_t Grid::map_centred(uint16_t index)
  {
    if (index == 0)
    {
      return centre;
    }

    if (index >= first_edge)
    {
      return map_centred_edge(index);
    }

    uint16_t ring = 1;
    int16_t ring_start = 1;
    int16_t ring_length = 8;
    while (ring_start + ring_length <= index)
    {
      ring++;
      ring_start += ring_length;
      ring_length += 8;
    }

    int16_t offset = index - ring_start;
    int16_t side = (offset << 2) / ring_length;

    offset -= (side * ring_length) >> 2;

    switch (side)
    {
    case 0:
      return centre + ring +
             (offset - ring + 1) * columns;

    case 1:
      return centre + ring - 1 - offset +
             ring * columns;

    case 2:
      return centre - ring -
             (offset - ring + 1) * columns;

    case 3:
      return centre - ring + 1 + offset -
             ring * columns;
    }
    return 0;
  }
} // namespace glow
