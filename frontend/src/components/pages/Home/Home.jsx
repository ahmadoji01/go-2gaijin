import React, { Component } from "react";
import Header from '../../elements/Header';
import HomeBanners from "../../elements/HomeBanners";
import ProductDisplaySection from "../../elements/ProductDisplaySection"
import './Home.scss';

class Home extends Component {
    constructor(props) {
        super(props);
        this.state = {
            chatHistory: []
        }
    }

    componentDidMount() {
        return fetch('http://localhost:8080/api/')
        .then((response) => response.json())
        .then((responseJson) => {
            console.log(responseJson);
        })
        .catch((error) => {
            console.error(error);
        });
    }

    render() {
        return (
            <div className="App">
                <Header />
                <HomeBanners />
                <ProductDisplaySection />
            </div>
        );    
    }
}
 
export default Home;