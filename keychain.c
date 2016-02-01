#include "keychain.h"

const char * cf_string_to_c_string_copy(CFStringRef str) {
  if (str == NULL) {
    return NULL;
  }

  CFIndex length = CFStringGetLength(str);
  CFIndex maxSize = CFStringGetMaximumSizeForEncoding(length, kCFStringEncodingUTF8) + 1;
  char *buffer = (char *)malloc(maxSize);
  if (CFStringGetCString(str, buffer, maxSize,
        kCFStringEncodingUTF8)) {
    return buffer;
  }
  return NULL;
}

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
  const char** services = (const char **) malloc(c * sizeof(char *));
  long idx = 0;
  for (i = 0; i < c; i++) {
    CFDictionaryRef dict = CFArrayGetValueAtIndex(result, i);
    if (CFDictionaryContainsKey(dict, CFSTR("desc"))) {
      continue;
    }
    CFStringRef svc = CFDictionaryGetValue(dict, CFSTR("svce"));
    services[idx++] = cf_string_to_c_string_copy(svc);
  }
  if (result != NULL) CFRelease(result);
  CFRelease(query);
  struct keychain_items *items = malloc(sizeof(struct keychain_items));
  items->Items = services;
  items->Count = idx;
  return items;
}

