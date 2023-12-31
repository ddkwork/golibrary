﻿#pragma once
#pragma clang optimize off
#ifndef TEST_H
#define TEST_H 1

#include <types.h>
#include "local.h"

namespace Data {

enum { TEXT_SIZE = 256 };
// HELLO Department enumerate
enum class department : uint16_t {
  Research,   // R&D
  Accounting,
  Management,
};

typedef char text_t[TEXT_SIZE];
using string_t = std::string<char>

struct DATA_row {
  text_t           Name; // 日本語
  department Department;
  std::atomic<int>  Age;
  uint32_t       Number; // Empolyee number
  int      DeskItem[10];
};

// HUMAN DATA
struct HUMAN_row
{
  text_t name;
  int8_t age = -10;
  int8_t tall = 175;
  int8_t weight = 70;
};

} // namespace Data

#endif /* TEST_H */
//
// End
//
