import { Loader } from 'pixi.js';
import React from 'react';
import './App.css';
import { Game } from './connection/socket';
import { GameState } from './game/state';
import Join from './join';
import Action from "./action"
import View from './view';

class App extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      game: {},
      credentials: {},
      joined: false,
      loader: Loader.shared,
      gameStarted: true,
      possibleActions: 0,
    }
  }

  componentDidMount() {
    this.state.loader.add("textures/cards.json");
  }

  possibleActionsChange(actions) {
    console.log("Action passthrough: ", actions)
    this.setState({ possibleActions: actions })
  }

  start(cred) {
    let gameState = new GameState();
    gameState.setOnPossibleActions(this.possibleActionsChange.bind(this))
    let game = new Game(gameState, cred, () => {
      this.setState({ joined: false });
    });
    game.start();
    this.setState({ game: game, joined: true });

    game.started.then(() => {
      console.log("Game started in app");
      this.setState({ gameStarted: true });
    });
  }

  onJoin(values) {
    this.start(values);
  }

  render() {
    return (
      <div className="App" >
        {
          this.state.joined
            ? <div>
              <Action game={this.state.game} actions={this.state.possibleActions}></Action>
              <View game={this.state.game} loader={this.state.loader}></View>
            </div>
            : <Join onJoin={this.onJoin.bind(this)}></Join>
        }
      </div>
    );
  }
}
export default App;
