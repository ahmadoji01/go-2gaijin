import React from "react";
import "./ProductDisplaySection.scss";
import { NavLink } from 'react-router-dom';
import ProductCard from '../ProductCard'

const ProductDisplaySection = () => (
    <section class="products">
        <div class="container">
            <div class="container">
                <div class="row">
                    <div class="col-md-12">
                        <div class="product-title-area">
                            <div class="product__title justify-content-center">
                                <div class="ml-auto align-content-center text-center">
                                    <h2>Recently Added</h2>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="row">
                    <ProductCard />
                </div>
                <div class="row">
                    <div class="col-md-12">
                        <div class="more-product">
                            <a href="all-products.html" class="btn btn--lg btn--round">All New Products</a>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </section>
);

export default ProductDisplaySection;