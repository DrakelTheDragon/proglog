# Protocol Buffers

Protocol buffers, also known as *protobuf*, is Google's language and platform-neutral extensible
mechanism for structuring and serializing data.

Protobuf lets you define how you want your data structured, compile your protobuf into code in
potentially many languages, and then read and write your structured data to and from different
data streams. Protocol buffers are good for communicating between two systems (such as
microservices), which is why Google used protobuf when building gRPC to develop a high-performance
remote procedure call (RPC) framework.

gRPC uses protocol buffers to define APIs and serialize messages.

## Advantages

The advantages of using protobuf are that it:

- Guarantees type-safety
- Prevents schema-violations
- Enables fast serialization
- Offers backward compatibility


## Example

Here's a qucik example that shows what protocol buffers look like and how they work.

Imagine you work at Twitter and one of the object types you work with are Tweets. Tweets, at the
very least, comprise the author's message. If you defined this in protobuf, it would look like this:

```proto
syntax = "proto3";

package twitter;

message Tweet {
    string message = 1;
}
```

You'd then compile this protobuf into code in the language of your choice. For example, the protobuf
compiler would take this protobuf and generate the following Go code:

```go
// Code generated by protoc-gen-go. DO NOT EDIT.
// source: example.proto

package twitter

type Tweet struct {
    Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

// Note: Protobuf generates internal fields and methods I haven't included for brevity.
```

## Why Use Protocol Buffers?

Protobuf offers all kinds of useful features:

### Consistent Schemas

With protobuf, you encode your semantics once and use them across your services to ensure a consistent
data model throughout your whole system.

### Versioning For Free

One of Google's motivations for creating protobuf was to eliminate the need for version checks and
prevent ugly code like this:

```go
if version == 3 {
    // ...
} else if version > 4 {
    if version == 5 {
        // ...
    }
    // ...
}
```

Think of a protobuf message like a Go struct because when you compile a message it turns into a struct.
With protobuf, you number your fields on your messages to ensure you maintain backward compatibility as
you roll out new features and changes to your protobuf. So it's easy to add new fields, and intermediate
servers that need not use the data can simply parse it and pass through it without needing to know about
all the fields. Likewise with removing fields: you can ensure that deprecated fields are no longer used
by marking them as reserved; the compiler will then complain if anyone tries to the deprecated fields.

### Less Boilerplate

The protobuf libraries handle encoding and decoding for you, which means you don't have to handwrite that
code yourself.

### Extensibility

The protobuf compiler supports extensions that can compile your protobuf into code using your own compilation
logic. For example, you might want several structs to have a common method. With protobuf, you can write a
plugin to generate that method automatically.

### Language Agnosticism

Protobuf is implemented in many languages: since Protobuf version 3.0, there's support for Go, C++, Java,
JavaScript, Python, Ruby, C#, Objective C, and PHP, and third-party support for other languages. And you
don't have to do any extra work to communicate between services written in different languages. This is
great for companies with various teams that want to use different languages, or when your team wants to
migrate to another language.

### Performance

Protobuf is highly performant, and has smaller payloads and serializes up to six times faster than JSON.