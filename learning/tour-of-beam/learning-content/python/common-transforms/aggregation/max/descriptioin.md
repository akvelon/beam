# Max

Max provides a variety of different transforms for computing the maximum values in a collection, either globally or for each key.
    
### Maximum element in a PCollection

You can use ```CombineGlobally(lambda elements: max(elements or [None]))``` to get the maximum element from the entire ```PCollection```.

```
import apache_beam as beam

with beam.Pipeline() as pipeline:
  max_element = (
      pipeline
      | 'Create numbers' >> beam.Create([3, 4, 1, 2])
      | 'Get max value' >>
      beam.CombineGlobally(lambda elements: max(elements or [None]))
      | beam.Map(print))
```

Output

```
4
```

### Maximum elements for each key

You can use ```Combine.PerKey()``` to get the maximum element for each unique key in a PCollection of key-values.

```
import apache_beam as beam

with beam.Pipeline() as pipeline:
  elements_with_max_value_per_key = (
      pipeline
      | 'Create produce' >> beam.Create([
          ('🥕', 3),
          ('🥕', 2),
          ('🍆', 1),
          ('🍅', 4),
          ('🍅', 5),
          ('🍅', 3),
      ])
      | 'Get max value per key' >> beam.CombinePerKey(max)
      | beam.Map(print))
```

Output

```
('🥕', 3)
('🍆', 1)
('🍅', 5)
```

### Description for example

Created a list of integers ```PCollection```. The ```beam.combiners.Top.Largest(5)``` to return the larger than [5] from `PCollection`.