## Dev environment tips
- Each time a data module, resource, or function is added in `internal/provider`, update the full example in `examples/full` so it can be demonstrated and tested.
- Update the tests in `internal/provider` any time a resource or data source is updated to ensure it is fully tested.

## Testing instructions
- Always run `make testacc` after code changes.
