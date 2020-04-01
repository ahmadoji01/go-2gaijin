import React from "react";
import "./ProductCard.scss";
import { NavLink } from 'react-router-dom';

const ProductCard = () => (
    <div className="col-lg-4 col-md-6">
        <div className="product product--card">
            <div className="product__thumbnail">
                <img src="images/p1.jpg" alt="Product Image" />
            </div>
            <div className="product-desc">
                <a href="single-product.html" className="product_title">
                    <h4>MartPlace Extension Bundle</h4>
                </a>
                <ul className="titlebtm">
                    <li>
                        <img className="auth-img" src="images/auth.jpg" alt="author image" />
                        <p>
                            <a href="#">AazzTech</a>
                        </p>
                    </li>
                    <li className="product_cat">
                        <a href="#">
                            <span className="lnr lnr-book"></span>Plugin</a>
                    </li>
                </ul>
                <p>Nunc placerat mi id nisi interdum mollis. Praesent pharetra, justo ut scelerisque the mattis,
                    leo quam aliquet congue.</p>
            </div>
            <div className="product-purchase">
                <div className="price_love">
                    <span>$10 - $50</span>
                    <p>
                        <span className="lnr lnr-heart"></span> 90</p>
                </div>
                <div className="sell">
                    <p>
                        <span className="lnr lnr-cart"></span>
                        <span>16</span>
                    </p>
                </div>
            </div>
        </div>
    </div>
);
    
export default ProductCard;