#include "keychain.h"

const char * cf_string_to_c_string_copy(CFStringRef str) {
  if (str == NULL) {
    return NULL;
  }

  CFIndex length = CFStringGetLength(str);
  CFIndex maxSize = CFStringGetMaximumSizeForEncoding(length, kCFStringEncodingUTF8) + 1;
  char *buffer = (char *)malloc(maxSize);
  if (CFStringGetCString(str, buffer, maxSize, kCFStringEncodingUTF8)) {
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
  const struct keychain_item** items = (const struct keychain_item **) malloc(c * sizeof(struct keychain_item*));
  long idx = 0;
  for (i = 0; i < c; i++) {
    CFDictionaryRef dict = CFArrayGetValueAtIndex(result, i);
    if (CFDictionaryContainsKey(dict, CFSTR("desc"))) {
      continue;
    }
    CFStringRef label = CFDictionaryGetValue(dict, CFSTR("labl"));
    CFStringRef svc = CFDictionaryGetValue(dict, CFSTR("svce"));
    CFStringRef acct = CFDictionaryGetValue(dict, CFSTR("acct"));
    struct keychain_item *item = malloc(sizeof(struct keychain_item));
    item->Label = cf_string_to_c_string_copy(label);
    item->Service = cf_string_to_c_string_copy(svc);
    item->Account = cf_string_to_c_string_copy(acct);
    items[idx++] = item;
  }
  if (result != NULL) {
    CFRelease(result);
  }
  CFRelease(query);
  struct keychain_items *items_struct = malloc(sizeof(struct keychain_items));
  items_struct->Items = items;
  items_struct->Count = idx;
  return items_struct;
}

void free_keychain_items(const struct keychain_items *itms) {
  if (itms == NULL) {
    return;
  }

  for (int i = 0; i < itms->Count; i++) {
    const struct keychain_item *itm = itms->Items[i];
    if (itm == NULL) {
      continue;
    }
    free((void *)itm->Account);
    free((void *)itm->Label);
    free((void *)itm->Service);
    free((void *)itm);
  }

  free(itms->Items);
  free((void *)itms);
}

