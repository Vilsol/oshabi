// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';
import {types} from '../models';

export function SetStream(arg1:boolean):Promise<Error>;

export function Calibrate():Promise<Error>;

export function GetConfig():Promise<main.ConvertedConfig>;

export function Read():Promise<Error>;

export function SetLeague(arg1:string):Promise<Error>;

export function UpdatePricing():Promise<Error>;

export function Copy():Promise<Error>;

export function GetDisplayCount():Promise<number>;

export function SetDisplay(arg1:number):Promise<Error>;

export function SetListingCount(arg1:string,arg2:number,arg3:number):Promise<Error>;

export function Clear():Promise<Error>;

export function SetName(arg1:string):Promise<Error>;

export function SetPrice(arg1:string,arg2:string):Promise<Error>;

export function GetListings():Promise<Array<types.ParsedListing>>;
