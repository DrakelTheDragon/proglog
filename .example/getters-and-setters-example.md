# Getters and Setters

## Getters

The protobuf compiler generates varius methods on structs, but the only methods you'll directly use are the getters. Use the struct's fields when you can,
but you'll find the getters useful when you have multiple messages with the same getter(s) and you want to abstract those method(s) into an interface. 

For example, imagine you're building a retail site like Amazon and have different types of stuff you sell - books, games, and so on - each with a field for
the item's price, and you want to find the total of the items in the user's cart. You'd make a `Pricer` interface and a `Total` function that takes in a
slice of `Pricer` interfaces and returns their total cost:

```go
type Book struct {
    Price uint64
}

func (b *Book) GetPrice() uint64 {
    // ...
}

type Game struct {
    Price uint64
}

func (g *Game) GetPrice() uint64 {
    // ...
}

type Pricer interface {
    GetPrice() uint64
}

func Total(items []Pricer) uint64 {
    // ...
}
```

## Setters

Now imagine that you want to write a script to change the price of all your inventory. You *could* do this with reflection, but reflection should be your
last resort since, as the Go proverb goes, reflection is never clear. If we just had setters, we could use an interface like the following to set the price
on the different kinds of items in your inventory:


```go
type PriceAdjuster interface {
    SetPrice(price uint64)
}
```

When the compiled code isn't quite what you need, you can extend the compiler's output with plugins.
