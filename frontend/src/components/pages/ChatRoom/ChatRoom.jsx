import React, { Component } from "react";
import Header from '../../elements/Header';
import ChatHistory from '../../elements/ChatHistory';
import ChatInput from '../../elements/ChatInput';
import { connect, sendMsg } from "../../../api";
import './ChatRoom.scss';

class ChatRoom extends Component {
    constructor(props) {
        super(props);
        this.state = {
            chatHistory: []
        }
    }

    send(event) {
        if(event.keyCode === 13 && event.target.value !== "") {
            sendMsg(event.target.value);
            event.target.value = "";
        }
    }

    componentDidMount() {
        connect((msg) => {
            console.log("New Message")
            this.setState(prevState => ({
            chatHistory: [...this.state.chatHistory, msg]
            }))
            console.log(this.state);
        });
    }

    render() {
        return (
            <div className="App">
                <Header />
                <ChatHistory chatHistory={this.state.chatHistory} />
                <ChatInput send={this.send} />
            </div>
        );
    }
}

export default ChatRoom;