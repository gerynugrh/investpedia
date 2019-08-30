import 'package:flutter/material.dart';
import 'package:sabang/home/home_page.dart';
import 'package:sabang/home/investment_page.dart';

class _HomeState extends State<Home> {
  int _currentIndex = 0;
  final List<Widget> _children = [HomePage(), InvestmentPage()];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Investpedia'),
      ),
      body: _children[_currentIndex],
      bottomNavigationBar: BottomNavigationBar(
        onTap: onTabTapped,
        currentIndex: _currentIndex,
        items: [
          BottomNavigationBarItem(
              icon: new Icon(Icons.home),
              title: new Text('Home')
          ),
          BottomNavigationBarItem(
              icon: new Icon(Icons.mail),
              title: new Text('Messages')
          ),
          BottomNavigationBarItem(
              icon: new Icon(Icons.person),
              title: new Text('Profile')
          )
        ],
      ),
    );
  }

  void onTabTapped(int index) {
    setState(() {
      _currentIndex = index;
    });
  }
}

class Home extends StatefulWidget {
  @override
  State<StatefulWidget> createState() {
    return _HomeState();
  }
}
