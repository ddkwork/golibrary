

static ZyanStatus ZydisInputPeek(ZydisDecoderContext *context,
                                 ZydisDecodedInstruction *instruction,
                                 ZyanU8 *value) {
  ZYAN_ASSERT(context);
  ZYAN_ASSERT(instruction);
  ZYAN_ASSERT(value);

  if (instruction->length >= ZYDIS_MAX_INSTRUCTION_LENGTH) {
    return ZYDIS_STATUS_INSTRUCTION_TOO_LONG;
  }

  if (context->buffer_len > 0) {
    *value = context->buffer[0];
    return ZYAN_STATUS_SUCCESS;
  }

  return ZYDIS_STATUS_NO_MORE_DATA;
}
