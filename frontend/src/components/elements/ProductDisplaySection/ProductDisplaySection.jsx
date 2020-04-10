import React, {Component} from 'react';
import "./ProductDisplaySection.scss";
import { NavLink } from 'react-router-dom';
import ProductCard from '../ProductCard'

class ProductDisplaySection extends Component {
    
    render() {
        if(typeof(this.props.items) !== 'undefined') {
            return (
                <section className="products">
                    <div className="container">
                        <div className="container">
                            <div className="row">
                                <div className="col-md-12">
                                    <div className="product-title-area">
                                        <div className="product__title justify-content-center">
                                            <div className="ml-auto align-content-center text-center">
                                                <h2>{this.props.title}</h2>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div className="row">
                            { this.props.items.map(function (item, i) {
                                return <ProductCard key={i} 
                                        name={item.name} 
                                        price={item.price} 
                                        sellerName={item.seller_name} 
                                        imgURL={item.img_url} 
                                        location={item.loc} />;
                            }
                            )}
                            </div>
                            <div className="row">
                                <div className="col-md-12">
                                    <div className="more-product">
                                        <a href="all-products.html" className="btn btn--lg btn--round">All New Products</a>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </section>
            );
        } else {
            return "";
        } 
    }
}

export default ProductDisplaySection;