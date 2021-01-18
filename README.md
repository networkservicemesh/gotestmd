# gotestmd

Tool to generate go tests based on markdown files.

## Usages

Generate suties with default bash runner:

```bash
gotestmd INPUT_DIR OUTPUT_DIR
```

Generate suites that using a custom runner:

```bash
gotestmd INPUT_DIR OUTPUT_DIR BASE_PKG
```


## Makrdown syntax

- `#Run` - **REQUIRED**  - Contains any text and `bash` steps. Can be any level, should be used once in a file. 
- `#Cleanup` - _OPTIONAL_ - Contains `bash` steps. Can be any level, should be used once in a file. 
- `#Requires` - _OPTIONAL_ - Contains a list of required dependencies in format markdown links.
- `#Includes` - _OPTIONAL_ -Contains a list of using examples in context of this example in format markdown links.

# Examples

See at [examples](./examples)
