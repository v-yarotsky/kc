#include "keychain.h"

const struct keychain_items * list_keychain_items() {
  CFMutableDictionaryRef query = CFDictionaryCreateMutable(kCFAllocatorDefault, 0,
                                    &kCFTypeDictionaryKeyCallBacks,
                                    &kCFTypeDictionaryValueCallBacks);

  CFDictionarySetValue(query, kSecReturnAttributes, kCFBooleanTrue);
  CFDictionarySetValue(query, kSecMatchLimit, kSecMatchLimitAll);
  CFDictionarySetValue(query, kSecClass, kSecClassGenericPassword);

  CFTypeRef result = NULL;
  SecItemCopyMatching(query, &result);
  CFIndex i, c = CFArrayGetCount(result);
  CFStringRef *services = (CFStringRef *) malloc(c * sizeof(CFStringRef));
  for (i = 0; i < c; i++) {
    CFDictionaryRef dict = CFArrayGetValueAtIndex(result, i);
    CFStringRef svc = CFDictionaryGetValue(dict, CFSTR("svce"));
    services[i] = svc;
  }
  /* if (result != NULL) CFRelease(result); */ // should really free this
  CFRelease(query);
  struct keychain_items *items = malloc(sizeof(struct keychain_items));
  items->Items = services;
  items->Count = c;
  return items;
}

