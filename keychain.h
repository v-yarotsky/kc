#ifndef KEYCHAIN_H
#define KEYCHAIN_H

#include <CoreFoundation/CoreFoundation.h>
#include <Security/SecKeychainItem.h>
#include <Security/SecItem.h>

struct keychain_item {
  const char *Label;
  const char *Service;
};

struct keychain_items {
  const struct keychain_item** Items;
  long Count;
};

#endif
