import React from "react";
import "./HomeBanners.scss";

const HomeBanners = () => (
    <section className="hero-area bgimage">
        <div className="hero-content content_above">
            <div className="content-wrapper">
                <div className="container">
                    <div className="row">
                        <div className="col-md-12">
                            <div className="container">
                                <div className="row">
                                    <div className="col-12">
                                        <div className="swiper-container">
                                            <div className="swiper-wrapper">
                                                <div className="swiper-slide">
                                                    <img className="img-fluid" src="https://via.placeholder.com/1368x768" alt="Slider 1" />
                                                </div>
                                                <div className="swiper-slide">
                                                    <img className="img-fluid" src="https://via.placeholder.com/1368x768" alt="Slider 2" />
                                                </div>
                                                <div className="swiper-slide">
                                                    <img className="img-fluid" src="https://via.placeholder.com/1368x768" alt="Slider 3" />
                                                </div>
                                                <div className="swiper-slide">
                                                    <img className="img-fluid" src="https://via.placeholder.com/1368x768" alt="Slider 4" />
                                                </div>
                                                <div className="swiper-slide">
                                                    <img className="img-fluid" src="https://via.placeholder.com/1368x768" alt="Slider 5" />
                                                </div>
                                            </div>
                                            <div className="swiper-pagination"></div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </section>
);

export default HomeBanners;