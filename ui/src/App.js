import React, { Component } from 'react';
import { BrowserRouter, Switch, Route } from 'react-router-dom';

import Navigation from './components/Navigation';
import Home from './components/Home';
import Label from './components/Label';

class App extends Component {
  render() {
      return (
        <BrowserRouter>
          <div>
            <Navigation />
            <Switch>
              <Route path="/" component={Home} exact={true} />
              <Route path="/label/:id" component={Label} />
            </Switch>
          </div>
        </BrowserRouter>
      );
  }
}

export default App;
