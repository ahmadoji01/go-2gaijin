import React from 'react';
import Header from '../../elements/Header';
import HomeBanners from "../../elements/HomeBanners";
import ProductDisplaySection from "../../elements/ProductDisplaySection"
import './Home.scss';
 
const home = () => {
    return (
        <div className="App">
            <Header />
            <HomeBanners />
            <ProductDisplaySection />
        </div>
    );
}
 
export default home;