import React from 'react';
import {
    Navbar,
    Nav,
    NavItem,
    NavLink } from 'reactstrap';

import { Link } from 'react-router-dom';
import moment from 'moment';
import i18n from 'meteor/universe:i18n';

const T = i18n.createComponent();

export default class Footer extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (
            <div>
                <Navbar color="light" light expand="md" fixed="bottom" id="footer" className="d-none d-md-flex">
                    <span className="text-muted"><a href="https://raw.githubusercontent.com/forbole/big_dipper/master/LICENSE" target="_blank"><T>navbar.license</T></a> &copy;2018-{moment().format('YYYY')}. </span>
                    <Nav className="ml-auto" navbar>
                        <NavItem>
                            <NavLink href="https://www.github.com/forbole/big_dipper" target="_blank"><i className="fab fa-github"></i> <T>navbar.forkMe</T></NavLink>
                        </NavItem>
                    </Nav>
                </Navbar>
                <Navbar color="light" light fixed="bottom" className="d-block d-md-none mobile-menu">
                    <Nav>
                        <NavItem>
                            <NavLink tag={Link} to="/"><i className="material-icons">home</i></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink tag={Link} to="/validators"><i className="material-icons">perm_contact_calendar</i></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink tag={Link} to="/blocks"><i className="fas fa-square"></i></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink tag={Link} to="/transactions"><i className="fas fa-sync"></i></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink tag={Link} to="/proposals"><i className="material-icons">insert_drive_file</i></NavLink>
                        </NavItem>
                        <NavItem>
                            <NavLink tag={Link} to="/voting-power-distribution"><i className="material-icons">power_on</i></NavLink>
                        </NavItem>
                    </Nav>
                </Navbar>
            </div>  
        );
    }
}
