import React from 'react';
import './App.css';
import View from './view';

class App extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      game: {
        active: false
      }
    }
  }

  componentDidMount() {
  }

  render() {
    return (
      <div className="App" >
        <View></View>
      </div>
    );
  }
}
export default App;
