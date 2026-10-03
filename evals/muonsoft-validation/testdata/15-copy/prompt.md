# Copy validation messages for a response

Implement CopyMessages: copy the existing map of field names to validation-message slices so response formatting can modify either map or any slice independently. Preserve nil maps and nil slice values, and do not rewrite field paths or message text. Keep the existing ValidateCount function unchanged. No new validation rules or framework are needed.
