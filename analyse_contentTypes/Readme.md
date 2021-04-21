# Content type report

Script that anaylse the Content model JSON generates by CLI or CMA and present a report on content types.

## requirements

go.1.14 up

## How to use

go run main.go -f some/file/fileName.json

### by executable
$ ./analyseContentTypes -f some/file/fileName.json

### planned

- executable application
- json file as parameter

## Example 

    $ ./analyseContentTypes -f example_app_content_model.json
    ** Analysis Report **
    Description:
    [Missing] - Value missing.
    [Notice] - Good practice attention.
    [Validation] - Validation missing or unsuported.

    Total ContentTypes: 10
    ***** REPORT ******* 

    ContentType name: Category
    ContentType id: category
    No  errors.

    ContentType name: Course
    ContentType id: course
    No  errors.

    ContentType name: Layout
    ContentType id: layout
    No  errors.

    ContentType name: Layout > Copy
    ContentType id: layoutCopy
    No  errors.

    ContentType name: Layout > Hero Image
    ContentType id: layoutHeroImage
    No  errors.

    ContentType name: Layout > Highlighted Course
    ContentType id: layoutHighlightedCourse
    No  errors.

    ContentType name: Lesson
    ContentType id: lesson
    No  errors.

    ContentType name: Lesson > Code Snippets
    ContentType id: lessonCodeSnippets
    No  errors.

    ContentType name: Lesson > Copy
    ContentType id: lessonCopy
    No  errors.

    ContentType name: Lesson > Image
    ContentType id: lessonImage
    No  errors.