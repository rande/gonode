export default class ActionDispatcher {
    private actions: Map<string, any> = new Map<string, any>();

    register(name: string, action: any): void {
        this.actions[name] = action;
    }

    dispatch(name: string, args: Array<any>): any | boolean {
        if (!(name in this.actions)) {
            console.log(`Action ${name} is not registered`);
            return false;
        }

        return this.actions[name](...args);
    }
}