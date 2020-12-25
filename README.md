# gotestmd

Tool to generate golang tests based on markdown examples

## Usages

```bash
gotestmd ${INPUTDIR} ${OUTPUTDIR}
```

## Example Syntax

`Run` **REQUIRED**  - Contains any text and `bash` steps. Can be any level, should be used once in a file. 

`Cleanup` **OPTIONAL** - Contains `bash` steps. Can be any level, should be used once in a file. 

`Requires` **OPTIONAL** - Contains a list of required dependencies in format markdown links.

`Includes` **OPTIONAL** -Contains a list of using examples in context of this example in format markdown links.


# Examples

See [examples](./examples)
