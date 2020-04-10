import React, { Component } from "react";
import Header from '../../elements/Header';
import HomeBanners from "../../elements/HomeBanners";
import ProductDisplaySection from "../../elements/ProductDisplaySection"
import './Home.scss';

class Home extends Component {
    
    constructor(props) {
        super(props);
        this.state = {
            data: []
        }
    }

    componentWillMount() {
        fetch('http://localhost:8080/api/')
        .then((response) => response.json())
        .then((responseJson) => {
            const jsonData = responseJson.data;
            this.setState({ data: jsonData});
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
                <ProductDisplaySection title="Featured Items" items={this.state.data.featureditems} />
            </div>
        );    
    }
}
 
export default Home;