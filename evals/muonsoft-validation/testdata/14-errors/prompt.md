# Format errors from the existing validation API

Implement ErrorText for the existing error types. Nil becomes an empty string. An error wrapping ErrInvalid returns "Invalid input". Other errors keep their Error() text. Preserve ValidateCode and its sentinel identity. This task changes error presentation only; retain the existing dependency-free validation implementation.
