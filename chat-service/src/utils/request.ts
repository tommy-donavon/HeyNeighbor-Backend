import { lookupService } from '../register/register.js';
import { ITenant, IProperty } from '../model';
import fetch from "node-fetch"

export const getUser = async (authHeader: string | undefined): Promise<ITenant | Error> => {
  if(typeof authHeader === undefined) return Promise.resolve(new Error('unauthorized'))
  try {
    const accountAddr = await lookupService('account-service');
    const response = await fetch(accountAddr, {
        headers:{
            'Authorization': authHeader as string,
        }
    })
    const user:ITenant = <ITenant> await response.json()
    return (response.status === 200 && user !== undefined ? Promise.resolve(user) : Promise.resolve(new Error('unable to process user')))
  } catch (err) {
    return Promise.resolve(new Error('unable to reach account-service'))
  }
};

export const getProperty = async (serverCode: string, authHeader:string): Promise<IProperty | Error> => {
    if(serverCode.length === 0) return Promise.resolve(new Error('invalid server code'))
    try {
        const propertyAddr = await lookupService('property-service')
        const response = await fetch(`${propertyAddr}code/${serverCode}`, {
            headers:{
                'Authorization': authHeader
            }
        })
        const property:IProperty = await response.json() as IProperty
        return (response.status === 200 && property ? Promise.resolve(property) : Promise.resolve(new Error('unable to retrieve property')))
    }catch(err){
        return Promise.resolve(err as Error)
    }
}