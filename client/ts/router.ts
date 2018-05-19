import Navigo from './navigo';

const root : any = null;
const useHash : boolean = true;
const hash : string = '#!';
const router = new Navigo(root, useHash, hash);

export default router;
