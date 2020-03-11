## Decimal-Go

This package contains a decimal type which works like the .Net decimal type and can be used whenever compatibility to this type is required
Basic conversion from most number types exist aswell as a conversion to float64

### Serialization

#### Read Decimal
```
value,err:=decimal.ReadDecimal(<io.Reader>)
```

Reads a decimal from a reader which points to a stream where a decimal in .Net format is stored.

#### Write Decimal
```
err := decimal.WriteDecimal(<io.Writer>, <decimal>)
```

Writes a decimal to a writer in a format which .Net understands

### Working with decimals

As arithmetic is not supported in the current state you should convert the decimal to some number type before you work with it and after you are done
convert it back to a decimal if you need to store it as a decimal