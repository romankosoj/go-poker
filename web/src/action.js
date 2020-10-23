import React from "react"

const { PLAYER_ACTION } = require("./events/constants")

class Action extends React.Component {
    constructor(props){
        super(props);
        this.game = props.game 
        this.state = {raise: 0};
        this.handleChange = this.handleChange.bind(this);
        this.fold = this.fold.bind(this);
        this.call = this.call.bind(this);
        this.check = this.check.bind(this);
        this.raise = this.raise.bind(this);
    }

    fold() {
        this.game.send({event: PLAYER_ACTION, data: {
            action: 1,
            payload: 0
        }})
    }

    call() {
        this.game.send({event: PLAYER_ACTION, data: {
            action: 2,
            payload: 0,
        }})
    }

    check(){
        this.game.send({event: PLAYER_ACTION, data: {
            action: 4,
            payload: 0,
        }})
    }

    raise(){
        this.game.send({event: PLAYER_ACTION, data: {
            action: 3,
            payload: this.state.raise,
        }})
    }

    handleChange(e){
        this.setState({raise: e.target.value});
    }

    render(){
        return (
            <div>
                <button onClick={this.fold}>FOLD</button>
                <button onClick={this.check}>CHECK</button>
                <button onClick={this.call}>CALL</button>
                <input type="number" id="raiseInput" name="raiseInput" value={this.state.raise} onChange={this.handleChange} />
                <button onClick={this.raise}>RAISE</button>
            </div>
        );
    }
}

export default Action