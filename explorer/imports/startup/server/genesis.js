import { Validators } from '../../api/validators/validators.js';
let validatorsCount = Validators.find({}).count();

if (validatorsCount == 0){
    console.log("no validators");
    
}