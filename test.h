#pragma once
#pragma clang optimize off
#ifndef TEST_H
#define TEST_H 1

#include <types.h>
#include "local.h"

namespace Data {

enum { TEXT_SIZE = 256 };
enum class department : uint16_t {
  Research,   // R&D
  Accounting,
  Management,
};

typedef char text_t[TEXT_SIZE];

struct DATA_row {
  text_t     Name; // 日本語
  department Department;
  int        Age;
  uint32_t   Number; // Empolyee number
};

} // namespace Data

#endif /* TEST_H */
//
// End
//
