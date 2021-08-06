# Script to copy content types between spaces

**IMPORTANT Before running:**
* This script is offer as it is with no responsability.
* All spaces must have the same locale (the file try to migrate locaes also).

## How to run

Usage: 

**$./run.sh [args]**

Options:
* -t  The Contentful management token to use
* -o  The Origin Contentful space
* -d  The Destination Contentful space
* -c  The list of types to migrate as csv
* -l  Option to migrate locales: yes/no.
* -h  Show help

**Example**:
```bash
$ ./run.sh -t CMA_TOKEN -o SPACE_ORIGIN -d SPACE_DESTINATION -c CSV_OF_CONTENTTYPES_IDS -l no
```




