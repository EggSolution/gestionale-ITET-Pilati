var ind;
var loghi;
var sez;
var loghiAlt;
var altezzaEff;
var impAperto;

window.onload = function(){
    ind = document.querySelector('.indicatore');
    loghi = document.querySelector('.iconeCentrali');
    sez = document.querySelector('.sezioni');
    imp = document.querySelector('.imp');

    impAperto = false;

    sez.style.transform = "translateY(-0%)";
    imp.style.opacity = "0";
    imp.style.zIndex = "0";

    loghiAlt = loghi.offsetTop;
    altezzaEff = loghiAlt - 63;
    ind.style.top = altezzaEff + "px";
}

function menu(sezione){
    switch(sezione){
        case 1:
            ind.style.top = altezzaEff + "px";
            sez.style.transform = "translateY(-0%)";
            imp.style.opacity = "0";
            imp.style.zIndex = "0";
            break;
        case 2:
            ind.style.top = altezzaEff + 54 + "px";
            sez.style.transform = "translateY(-25%)";
            imp.style.opacity = "0";
            imp.style.zIndex = "0";
            break;
        case 3:
            ind.style.top = altezzaEff + 54 * 2 + "px";
            sez.style.transform = "translateY(-50%)";
            imp.style.opacity = "0";
            imp.style.zIndex = "0";
            break;
        case 4:
            ind.style.top = altezzaEff + 54 * 3 + "px";
            sez.style.transform = "translateY(-75%)";
            imp.style.opacity = "0";
            imp.style.zIndex = "0";
            break;
        case 5:
            if(impAperto == true){
                impAperto = false;
                imp.style.opacity = "0";
                imp.style.zIndex = "0";
            } else if(impAperto == false){
                impAperto = true;
                imp.style.opacity = "1";
                imp.style.zIndex = "5";
            }
            break;
        default:
            alert("errore nello scorrimento delle sezioni");
    }
}