let mobile = styles => Css.media("(max-width: 768px)", styles);
let smallMobile = styles => Css.media("(max-width: 370px)", styles);

let isMobile = () => {
  Webapi.Dom.(window |> Window.innerWidth <= 768);
};

let isSmallMobile = () => {
  Webapi.Dom.(window |> Window.innerWidth <= 370);
};
