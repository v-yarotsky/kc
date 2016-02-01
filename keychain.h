#ifndef KEYCHAIN_H
#define KEYCHAIN_H

#include <CoreFoundation/CoreFoundation.h>
#include <Security/SecKeychainItem.h>
#include <Security/SecItem.h>

struct keychain_items {
  CFStringRef * Items;
  long Count;
};

#endif
