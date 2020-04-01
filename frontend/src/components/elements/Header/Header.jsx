import React from "react";
import "./Header.scss";
import { NavLink } from 'react-router-dom';

const Header = () => (
  <div className="menu-area">
      <div className="topbar">
          <div className="topbar__container container">
              <div className="topbar__row">
                  <div className="topbar__item topbar__item--link">
                      <a className="topbar-link" href="#">About Us</a>
                  </div>
                  <div className="topbar__item topbar__item--link">
                      <a className="topbar-link" href="#">Contact Us</a>
                  </div>
                  <div className="topbar__item topbar__item--link">
                      <a className="topbar-link" href="#">Terms and Policies</a>
                  </div>
                  <div className="topbar__item topbar__item--link">
                      <a className="topbar-link" href="#">FAQ</a>
                  </div>
                  <div className="topbar__spring"></div>
                  <div className="topbar__item">
                      <div className="topbar-dropdown">
                          <button className="topbar-dropdown__btn" type="button">Currency:
                              <span className="topbar__item-value">JPY (&yen;)</span>
                          </button>
                      </div>
                  </div>
                  <div className="topbar__item">
                      <div className="topbar-dropdown">
                          <button className="topbar-dropdown__btn" type="button">Language:
                              <span className="topbar__item-value">EN</span>
                          </button>
                      </div>
                  </div>
              </div>
          </div>
      </div>
      <div className="top-menu-area">
          <div className="container">
              <div className="row">
                  <div className="col-lg-3 col-md-3 col-6 v_middle">
                      <div className="logo">
                          <NavLink to="/">
                              <img src="images/logo.png" alt="logo image" className="img-fluid" />
                          </NavLink>
                      </div>
                  </div>
                  <div className="col-lg-8 offset-lg-1 col-md-9 col-6 v_middle">
                      <div className="author-area">
                          <div className="join inline">
                              <a href="signup.html" className="btn author-area__seller-btn btn--xs">Become a Seller</a>
                          </div>
                          <div className="author__notification_area">
                              <ul>
                                  <li className="has_dropdown">
                                      <div className="icon_wrap">
                                          <span className="lnr lnr-alarm"></span>
                                          <span className="notification_count noti">25</span>
                                      </div>
                                      <div className="dropdowns notification--dropdown">
                                          <div className="dropdown_module_header">
                                              <h4>Notifications</h4>
                                              <a href="notification.html">View All</a>
                                          </div>
                                          <div className="notifications_module">
                                              <div className="notification">
                                                  <div className="notification__info">
                                                      <div className="info_avatar">
                                                          <img src="images/notification_head.png" alt="" />
                                                      </div>
                                                      <div className="info">
                                                          <p>
                                                              <span>Anderson</span> added to Favourite
                                                              <a href="#">Mccarther Coffee Shop</a>
                                                          </p>
                                                          <p className="time">Just now</p>
                                                      </div>
                                                  </div>
                                                  <div className="notification__icons ">
                                                      <span className="lnr lnr-heart loved noti_icon"></span>
                                                  </div>
                                              </div>
                                              <div className="notification">
                                                  <div className="notification__info">
                                                      <div className="info_avatar">
                                                          <img src="images/notification_head2.png" alt="" />
                                                      </div>
                                                      <div className="info">
                                                          <p>
                                                              <span>Michael</span> commented on
                                                              <a href="#">MartPlace Extension Bundle</a>
                                                          </p>
                                                          <p className="time">Just now</p>
                                                      </div>
                                                  </div>
                                                  <div className="notification__icons ">
                                                      <span className="lnr lnr-bubble commented noti_icon"></span>
                                                  </div>
                                              </div>
                                              <div className="notification">
                                                  <div className="notification__info">
                                                      <div className="info_avatar">
                                                          <img src="images/notification_head3.png" alt="" />
                                                      </div>
                                                      <div className="info">
                                                          <p>
                                                              <span>Khamoka </span>purchased
                                                              <a href="#"> Visibility Manager Subscriptions</a>
                                                          </p>
                                                          <p className="time">Just now</p>
                                                      </div>
                                                  </div>
                                                  <div className="notification__icons ">
                                                      <span className="lnr lnr-cart purchased noti_icon"></span>
                                                  </div>
                                              </div>
                                              <div className="notification">
                                                  <div className="notification__info">
                                                      <div className="info_avatar">
                                                          <img src="images/notification_head4.png" alt="" />
                                                      </div>
                                                      <div className="info">
                                                          <p>
                                                              <span>Anderson</span> added to Favourite
                                                              <a href="#">Mccarther Coffee Shop</a>
                                                          </p>
                                                          <p className="time">Just now</p>
                                                      </div>
                                                  </div>
                                                  <div className="notification__icons ">
                                                      <span className="lnr lnr-star reviewed noti_icon"></span>
                                                  </div>
                                              </div>
                                          </div>
                                      </div>
                                  </li>
                                  <li className="has_dropdown">
                                      <div className="icon_wrap">
                                          <span className="lnr lnr-envelope"></span>
                                          <span className="notification_count msg">6</span>
                                      </div>
                                      <div className="dropdowns messaging--dropdown">
                                          <div className="dropdown_module_header">
                                              <h4>My Messages</h4>
                                              <NavLink to="/chatroom">View All</NavLink>
                                          </div>
                                          <div className="messages">
                                              <a href="message.html" className="message recent">
                                                  <div className="message__actions_avatar">
                                                      <div className="avatar">
                                                          <img src="images/notification_head4.png" alt="" />
                                                      </div>
                                                  </div>
                                                  <div className="message_data">
                                                      <div className="name_time">
                                                          <div className="name">
                                                              <p>NukeThemes</p>
                                                              <span className="lnr lnr-envelope"></span>
                                                          </div>
                                                          <span className="time">Just now</span>
                                                          <p>Hello John Smith! Nunc placerat mi ...</p>
                                                      </div>
                                                  </div>
                                              </a>
                                              <a href="message.html" className="message recent">
                                                  <div className="message__actions_avatar">
                                                      <div className="avatar">
                                                          <img src="images/notification_head5.png" alt="" />
                                                      </div>
                                                  </div>
                                                  <div className="message_data">
                                                      <div className="name_time">
                                                          <div className="name">
                                                              <p>Crazy Coder</p>
                                                              <span className="lnr lnr-envelope"></span>
                                                          </div>
                                                          <span className="time">Just now</span>
                                                          <p>Hi! Nunc placerat mi id nisi interum ...</p>
                                                      </div>
                                                  </div>
                                              </a>
                                              <a href="message.html" className="message">
                                                  <div className="message__actions_avatar">
                                                      <div className="avatar">
                                                          <img src="images/notification_head6.png" alt="" />
                                                      </div>
                                                  </div>
                                                  <div className="message_data">
                                                      <div className="name_time">
                                                          <div className="name">
                                                              <p>Hybrid Insane</p>
                                                          </div>
                                                          <span className="time">Just now</span>
                                                          <p>Hi! Nunc placerat mi id nisi interum ...</p>
                                                      </div>
                                                  </div>
                                              </a>
                                              <a href="message.html" className="message">
                                                  <div className="message__actions_avatar">
                                                      <div className="avatar">
                                                          <img src="images/notification_head3.png" alt="" />
                                                      </div>
                                                  </div>
                                                  <div className="message_data">
                                                      <div className="name_time">
                                                          <div className="name">
                                                              <p>ThemeXylum</p>
                                                          </div>
                                                          <span className="time">Just now</span>
                                                          <p>Hi! Nunc placerat mi id nisi interum ...</p>
                                                      </div>
                                                  </div>
                                              </a>
                                              <a href="message.html" className="message">
                                                  <div className="message__actions_avatar">
                                                      <div className="avatar">
                                                          <img src="images/notification_head4.png" alt="" />
                                                      </div>
                                                  </div>
                                                  <div className="message_data">
                                                      <div className="name_time">
                                                          <div className="name">
                                                              <p>NukeThemes</p>
                                                              <span className="lnr lnr-envelope"></span>
                                                          </div>
                                                          <span className="time">Just now</span>
                                                          <p>Hello John Smith! Nunc placerat mi ...</p>
                                                      </div>
                                                  </div>
                                              </a>
                                          </div>
                                      </div>
                                  </li>
                                  <li className="has_dropdown">
                                      <div className="icon_wrap">
                                          <span className="lnr lnr-cart"></span>
                                          <span className="notification_count purch">2</span>
                                      </div>
                                      <div className="dropdowns dropdown--cart">
                                          <div className="cart_area">
                                              <div className="cart_product">
                                                  <div className="product__info">
                                                      <div className="thumbn">
                                                          <img src="images/capro1.jpg" alt="cart product thumbnail" />
                                                      </div>
                                                      <div className="info">
                                                          <a className="title" href="single-product.html">Finance and Consulting Business Theme</a>
                                                          <div className="cat">
                                                              <a href="#">
                                                                  <img src="images/catword.png" alt="" />Wordpress</a>
                                                          </div>
                                                      </div>
                                                  </div>
                                                  <div className="product__action">
                                                      <a href="#">
                                                          <span className="lnr lnr-trash"></span>
                                                      </a>
                                                      <p>$60</p>
                                                  </div>
                                              </div>
                                              <div className="cart_product">
                                                  <div className="product__info">
                                                      <div className="thumbn">
                                                          <img src="images/capro2.jpg" alt="cart product thumbnail" />
                                                      </div>
                                                      <div className="info">
                                                          <a className="title" href="single-product.html">Flounce - Multipurpose OpenCart Theme</a>
                                                          <div className="cat">
                                                              <a href="#">
                                                                  <img src="images/catword.png" alt="" />Wordpress</a>
                                                          </div>
                                                      </div>
                                                  </div>
                                                  <div className="product__action">
                                                      <a href="#">
                                                          <span className="lnr lnr-trash"></span>
                                                      </a>
                                                      <p>$60</p>
                                                  </div>
                                              </div>
                                              <div className="total">
                                                  <p>
                                                      <span>Total :</span>$80</p>
                                              </div>
                                              <div className="cart_action">
                                                  <a className="go_cart" href="cart.html">View Cart</a>
                                                  <a className="go_checkout" href="checkout.html">Checkout</a>
                                              </div>
                                          </div>
                                      </div>
                                  </li>
                              </ul>
                          </div>
                          <div className="author-author__info inline has_dropdown">
                              <div className="author__avatar">
                                  <img src="images/usr_avatar.png" alt="user avatar" />
                              </div>
                              <div className="autor__info">
                                  <p className="name">
                                      Jhon Doe
                                  </p>
                                  <p className="ammount">&yen;20.45</p>
                              </div>
                              <div className="dropdowns dropdown--author">
                                  <ul>
                                      <li>
                                          <a href="author.html">
                                              <span className="lnr lnr-user"></span>Profile</a>
                                      </li>
                                      <li>
                                          <a href="dashboard.html">
                                              <span className="lnr lnr-home"></span> Dashboard</a>
                                      </li>
                                      <li>
                                          <a href="dashboard-setting.html">
                                              <span className="lnr lnr-cog"></span> Setting</a>
                                      </li>
                                      <li>
                                          <a href="cart.html">
                                              <span className="lnr lnr-cart"></span>Purchases</a>
                                      </li>
                                      <li>
                                          <a href="favourites.html">
                                              <span className="lnr lnr-heart"></span> Favourite</a>
                                      </li>
                                      <li>
                                          <a href="dashboard-add-credit.html">
                                              <span className="lnr lnr-dice"></span>Add Credits</a>
                                      </li>
                                      <li>
                                          <a href="dashboard-statement.html">
                                              <span className="lnr lnr-chart-bars"></span>Sale Statement</a>
                                      </li>
                                      <li>
                                          <a href="dashboard-upload.html">
                                              <span className="lnr lnr-upload"></span>Upload Item</a>
                                      </li>
                                      <li>
                                          <a href="dashboard-manage-item.html">
                                              <span className="lnr lnr-book"></span>Manage Item</a>
                                      </li>
                                      <li>
                                          <a href="dashboard-withdrawal.html">
                                              <span className="lnr lnr-briefcase"></span>Withdrawals</a>
                                      </li>
                                      <li>
                                          <a href="#">
                                              <span className="lnr lnr-exit"></span>Logout</a>
                                      </li>
                                  </ul>
                              </div>
                          </div>
                      </div>
                      <div className="mobile_content ">
                          <span className="lnr lnr-user menu_icon"></span>
                          <div className="offcanvas-menu closed">
                              <span className="lnr lnr-cross close_menu"></span>
                              <div className="author-author__info">
                                  <div className="author__avatar v_middle">
                                      <img src="images/usr_avatar.png" alt="user avatar" />
                                  </div>
                                  <div className="autor__info v_middle">
                                      <p className="name">
                                          Jhon Doe
                                      </p>
                                      <p className="ammount">$20.45</p>
                                  </div>
                              </div>
                              <div className="author__notification_area">
                                  <ul>
                                      <li>
                                          <a href="notification.html">
                                              <div className="icon_wrap">
                                                  <span className="lnr lnr-alarm"></span>
                                                  <span className="notification_count noti">25</span>
                                              </div>
                                          </a>
                                      </li>
                                      <li>
                                          <a href="message.html">
                                              <div className="icon_wrap">
                                                  <span className="lnr lnr-envelope"></span>
                                                  <span className="notification_count msg">6</span>
                                              </div>
                                          </a>
                                      </li>
                                      <li>
                                          <a href="cart.html">
                                              <div className="icon_wrap">
                                                  <span className="lnr lnr-cart"></span>
                                                  <span className="notification_count purch">2</span>
                                              </div>
                                          </a>
                                      </li>
                                  </ul>
                              </div>
                              <div className="dropdowns dropdown--author">
                                  <ul>
                                      <li>
                                          <a href="author.html">
                                              <span className="lnr lnr-user"></span>Profile</a>
                                      </li>
                                      <li>
                                          <a href="dashboard.html">
                                              <span className="lnr lnr-home"></span> Dashboard</a>
                                      </li>
                                      <li>
                                          <a href="dashboard-setting.html">
                                              <span className="lnr lnr-cog"></span> Setting</a>
                                      </li>
                                      <li>
                                          <a href="cart.html">
                                              <span className="lnr lnr-cart"></span>Purchases</a>
                                      </li>
                                      <li>
                                          <a href="favourites.html">
                                              <span className="lnr lnr-heart"></span> Favourite</a>
                                      </li>
                                      <li>
                                          <a href="dashboard-add-credit.html">
                                              <span className="lnr lnr-dice"></span>Add Credits</a>
                                      </li>
                                      <li>
                                          <a href="dashboard-statement.html">
                                              <span className="lnr lnr-chart-bars"></span>Sale Statement</a>
                                      </li>
                                      <li>
                                          <a href="dashboard-upload.html">
                                              <span className="lnr lnr-upload"></span>Upload Item</a>
                                      </li>
                                      <li>
                                          <a href="dashboard-manage-item.html">
                                              <span className="lnr lnr-book"></span>Manage Item</a>
                                      </li>
                                      <li>
                                          <a href="dashboard-withdrawal.html">
                                              <span className="lnr lnr-briefcase"></span>Withdrawals</a>
                                      </li>
                                      <li>
                                          <a href="#">
                                              <span className="lnr lnr-exit"></span>Logout</a>
                                      </li>
                                  </ul>
                              </div>
                          </div>
                      </div>
                  </div>
              </div>
          </div>
      </div>
      <div className="mainmenu">
          <div className="container">
              <div className="row">
                  <div className="col-lg-8 col-md-8 col-sm-8">
                      <div className="navbar-header">
                          <nav className="navbar navbar-expand-md navbar-light mainmenu__menu">
                              <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
                                  <span className="navbar-toggler-icon"></span>
                              </button>
                              <div className="collapse navbar-collapse" id="navbarNav">
                                  <ul className="navbar-nav">
                                      <li className="has_megamenu">
                                          <a href="#"><i className="fa fa-fw fa-list"></i>Categories</a>
                                          <div className="dropdown_megamenu">
                                              <div className="megamnu_module">
                                                  <div className="menu_items">
                                                      <div className="menu_column">
                                                          <ul>
                                                              <li className="title"><i className="fa fa-fw fa-plug"></i> Electronics</li>
                                                              <li>
                                                                  <a href="#">Test 1</a>
                                                              </li>
                                                              <li>
                                                                  <a href="#">Test 2</a>
                                                              </li>
                                                              <li>
                                                                  <a href="#">Test 3</a>
                                                              </li>
                                                              <li>
                                                                  <a href="#">Test 4</a>
                                                              </li>
                                                              <li className="title"><i className="fa fa-fw fa-headphones"></i> Audio</li>
                                                              <li>
                                                                  <a href="#">Test 1</a>
                                                              </li>
                                                              <li>
                                                                  <a href="#">Test 2</a>
                                                              </li>
                                                          </ul>
                                                      </div>
                                                      <div className="menu_column">
                                                          <ul>
                                                              <li className="title"><i className="fa fa-fw fa-chair"></i> Furnitures</li>
                                                              <li>
                                                                  <a href="#">Test 1</a>
                                                              </li>
                                                              <li>
                                                                  <a href="#">Test 2</a>
                                                              </li>
                                                              <li className="title"><i className="fa fa-fw fa-utensils"></i> Kitchens</li>
                                                              <li>
                                                                  <a href="#">Test 1</a>
                                                              </li>
                                                              <li>
                                                                  <a href="#">Test 2</a>
                                                              </li>
                                                          </ul>
                                                      </div>
                                                      <div className="menu_column">
                                                          <ul>
                                                              <li className="title"><i className="fa fa-fw fa-car"></i> Cars</li>
                                                              <li>
                                                                  <a href="#">Accessories</a>
                                                              </li>
                                                              <li>
                                                                  <a href="#">Car Audio</a>
                                                              </li>
                                                              <li>
                                                                  <a href="#">Test 3</a>
                                                              </li>
                                                              <li className="title"><i className="fa fa-fw fa-motorcycle"></i> Motorcycle</li>
                                                              <li>
                                                                  <a href="#">Helm</a>
                                                              </li>
                                                              <li>
                                                                  <a href="#">Accessories</a>
                                                              </li>
                                                          </ul>
                                                      </div>
                                                  </div>
                                              </div>
                                          </div>
                                      </li>
                                      <li>
                                          <a href="#"><i className="fa fa-fw fa-fire"></i> Free Products</a>
                                      </li>
                                  </ul>
                              </div>
                          </nav>
                      </div>
                  </div>
                  <div className="col-lg-4 col-md-4 col-sm-4">
                      <div className="mainmenu__search">
                          <form action="#">
                              <div className="searc-wrap">
                                  <div className="input-group">
                                      <input type="text" placeholder="Find Something" />
                                      <div className="input-group-append">
                                          <button className="btn btn-primary btn-group-lg search-wrap__btn" type="submit">
                                              <span className="lnr lnr-magnifier"></span>
                                          </button>
                                      </div>
                                  </div>
                              </div>
                          </form>
                      </div>
                  </div>
              </div>
          </div>
      </div>
  </div>
);

export default Header;