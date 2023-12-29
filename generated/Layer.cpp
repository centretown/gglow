#include "Layer.h"

namespace glow
{
#ifndef MICRO_CONTROLLER
  std::string Layer::make_code()
  {
    std::stringstream s;
    s << "{" << length << ","
      << rows << ","
      << grid.make_code() << ","
      << chroma.make_code() << ","
      << hue_shift << ","
      << scan << ","
      << begin << ","
      << end << "}";
    return s.str();
  }

  std::string Layer::keys[Layer::KEY_COUNT] = {
      "length",
      "rows",
      "grid",
      "chroma",
      "hue_shift",
      "scan",
      "begin",
      "end",
  };
#endif

  void Layer::set_bounds()
  {
    auto ratio = [](uint16_t offset, uint16_t length)
    {
      if (offset > 100)
        offset %= 100;
      return static_cast<float>(offset) / 100.0f *
             static_cast<float>(length);
    };

    first = grid.adjust_bounds(ratio(begin, length));
    last = grid.adjust_bounds(ratio(end, length));

    if (last < first)
    {
      std::swap(first, last);
    }
  }

}