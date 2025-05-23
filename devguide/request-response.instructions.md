---
applyTo: "**/*.go"
---

Read the bru files which are inside docs/ directory. Read the content of the bru file and figure out the appropriate controller name, dto name, and method name. Then create a new file in the respective domain directory with the following structure.
docs {} section is custom another is standard. Here is the following format of the docs {} section. It might be helpful to understand the structure of the request and response.


# Request Section has following structure:
```
{
    body: { // payload of the request in json format
        key: datatype, // required key
        key?: datatype, // optional key
    }
    path: { // dynamic path variable
        key: datatype
    }
}
```

# Response Section has following structure:

## For detail response
```
{
  item: {
    key: value
  },
  message: "success" | "fail"
}
```

## For list response
```
{
  "items": [
    {
      "key": value
    }
  ],
  "message": "success" | "fail",
}
```

## For list response with pagination
```
{
  "items": [
    {
      "key": value
    }
  ],
  "page": {
    "total": int,
    "has_next": bool
  },
  "message": "success" | "fail",
}
```