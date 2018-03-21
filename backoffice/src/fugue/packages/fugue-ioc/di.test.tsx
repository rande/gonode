
import { Container, Definition, ContainerBuilder } from './di';

class FooBar {
    private counter: number;

    constructor(counter: number) {
        this.counter = counter;
    }
}

function ServiceConstructor(arg: string) {
    return () => {
        return arg;
    };
}

function AnyService(...args: Array<any>) {
    return () => {
        return args;
    };
}

describe('Container', () => {
    it('should get a value', () => {
        let c = new Container();

        c.set('foo', 'bar');
        let [r, err] = c.get('foo');

        expect(r).toBe('bar');
        expect(err).toBe(null);
    });

    it('test container', () => {
        let c = new Container();
        let [s, err] = c.get('non existant service');

        expect(s).toBe(null);
        expect(err).toBeDefined();
    });

    it('test service builder', () => {
        let b = new ContainerBuilder();
        let c = new Container();

        b.set('service.constructor', new Definition(ServiceConstructor, ["hello thomas"]));
        b.set('foo.42', new Definition(FooBar, [42]));
        b.set('foo.2', new Definition(FooBar, [2]));
        b.set('foo.1', new Definition(FooBar, [1]));

        b.build(c);

        let [s1, err1] = c.get('service.constructor');

        expect(s1).toBeDefined();
        expect(err1).toBe(null);

        let [s2, err2] = c.get('foo.42');
        expect(s2).toBeDefined();
        expect(err2).toBe(null);
        expect(s2.counter).toBe(42);
    });

    it('test service alias', () => {
        let b = new ContainerBuilder();
        let c = new Container();

        b.def('service.child', AnyService, ['child']);
        b.def('service.parent', AnyService, ['parent', '@service.child']);

        b.build(c);

        let [s, err] = c.get('service.parent');
        expect(s).toBeDefined();
        expect(err).toBe(null);

        let [arg1, arg2] = s();

        expect(arg1).toEqual('parent');
        expect(typeof arg2 === 'function').toBe(true);
        expect(arg2()).toEqual(['child']);
    });

    it('test infinite loop with circular reference', () => {
        let b = new ContainerBuilder();
        let c = new Container();

        b.def('service.child', AnyService, ['child', '@service.parent']);
        b.def('service.parent', AnyService, ['parent', '@service.child']);

        let [container, err] = b.build(c);

        expect(container.services).toBeDefined();
        expect(err).toBeDefined();
    });
});