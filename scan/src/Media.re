let mobile = (styles) => Css.media("(max-width: 576px)", styles);

let isMobile = () => {
  Window.innerWidth <= 576;
};
