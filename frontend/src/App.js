import React, { Component } from "react";
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import "./App.css";
import Home from './components/pages/Home';
import ChatRoom from './components/pages/ChatRoom';
import Error from './components/pages/Error'

class App extends Component {

  render() {
    return (      
      <BrowserRouter>
        <div>
          <Switch>
            <Route path="/" component={Home} exact/>
            <Route path="/chatroom" component={ChatRoom}/>
            <Route path="/product" component={ChatRoom}/>
          <Route component={Error}/>
          </Switch>
        </div> 
      </BrowserRouter>
   );
  }

}

export default App;