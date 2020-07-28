let mobile = styles => Css.media("(max-width: 768px)", styles);

let isMobile = () => {
  Window.innerWidth <= 768;
};
