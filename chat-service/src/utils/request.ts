import { lookupService } from '../register/register.js';
import { ITenant, IProperty } from '../model';
import fetch from "node-fetch"

export const getUser = async (authHeader: string | undefined): Promise<ITenant> => {
  if(typeof authHeader === undefined) throw new Error('missing auth header')
  try {
    const accountAddr = await lookupService('account-service');
    const response = await fetch(accountAddr, {
        headers:{
            'Authorization': authHeader as string,
        }
    })
    const user:ITenant = <ITenant> await response.json()
    if (response.ok && user) {
        return Promise.resolve(user)
    }
    throw new Error('unable to retrieve user')
  } catch (err) {
    throw err
  }
};

export const getProperty = async (serverCode: string, authHeader:string): Promise<IProperty> => {
    if(serverCode.length === 0)throw new Error('invalid server code')
    try {
        const propertyAddr = await lookupService('property-service')
        const response = await fetch(`${propertyAddr}code/${serverCode}`, {
            headers:{
                'Authorization': authHeader
            }
        })
        const property:IProperty = await response.json() as IProperty
        if(response.ok && property){
            Promise.resolve(property)
        }
        throw new Error('unable to retrieve property')
    }catch(err){
        throw err
    }
}