import 'package:playground/modules/examples/models/example_model.dart';

class CategoryModel {
  final String name;
  final List<ExampleModel> examples;

  CategoryModel(this.name, this.examples);
}