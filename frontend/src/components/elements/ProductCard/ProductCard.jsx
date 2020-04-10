import React, { Component } from 'react';
import "./ProductCard.scss";
import { NavLink } from 'react-router-dom';

class ProductCard extends Component {
    
    render() {
        return(
            <div className="col-lg-3 col-md-4">
                <div className="product product--card">
                    <div className="product__thumbnail">
                        <img src={ this.props.imgURL } alt="Product Image" />
                    </div>
                    <div className="product-desc">
                        <a href="single-product.html" className="product_title">
                            <h4>{ this.props.name }</h4>
                        </a>
                        <ul className="titlebtm">
                            <li>
                                <img className="auth-img" src="images/author-avatar.jpg" alt="author image" />
                                <p>
                                    <a href="#">{ this.props.sellerName }</a>
                                </p>
                            </li>
                            <li className="product_cat">
                                <a href="#">
                                    <span className="lnr lnr-book"></span>Plugin</a>
                            </li>
                        </ul>
                    </div>
                    <div className="product-purchase">
                        <div className="price_love">
                            <span>{ this.props.price }</span>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}
    
export default ProductCard;