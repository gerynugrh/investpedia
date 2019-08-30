import 'package:flutter/material.dart';

class HomePage extends StatefulWidget {
  @override
  State createState() {
    return _HomePageState();
  }
}

class _HomePageState extends State<HomePage> {

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Text('Home Page')
    );
  }
}