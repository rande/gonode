import {Email, MinLength, MaxLength, NotEmpty, Url, MustContains} from './index';

const tests = [
    ['Test Email: invalid emails',  Email, [], 'thomas.rabaix@gmail', {"keys": {}, "msg": "invalid email"}],
    ['Test Email: valid emails',  Email, [], 'thomas.rabaix@gmail.com', null],

    ['Test MinLength: valid string',            MinLength, [2], 'hello', null],
    ['Test MinLength: invalid length',          MinLength, [10], 'hello', {"keys": {"length": "5", "min": 10}, "msg": "invalid length"}],
    ['Test MinLength: invalid string (number)', MinLength, [10], 18, {"keys": {}, "msg": "invalid string"}],

    ['Test MaxLength: string',          MaxLength, [2], 'hello', {"keys": {"length": "5", "max": 2}, "msg": "invalid length"}],
    ['Test MaxLength: invalid length',  MaxLength, [10], 'hello', null],
    ['Test MaxLength: invalid length',  MaxLength, [10], 18, {"keys": {}, "msg": "invalid string"}],

    ['Test NotEmpty: empty string', NotEmpty, [], '', {"keys": {}, "msg": "value is empty"}],
    ['Test NotEmpty: empty object', NotEmpty, [], {}, {"keys": {}, "msg": "value is empty"}],
    ['Test NotEmpty: empty array', NotEmpty, [], [], {"keys": {}, "msg": "value is empty"}],
    ['Test NotEmpty: valid string', NotEmpty, [], 'hello', null],
    ['Test NotEmpty: valid object', NotEmpty, [], {foo: 2}, null],
    ['Test NotEmpty: valid array', NotEmpty, [], [1], null],

    ['Test Url: invalid Url',  Url, [], 'http:foo.bar', {"keys": {}, "msg": "invalid url"}],
    ['Test Url: valid Url',  Url, [], 'http://foo.bar', null],

    ['Test MustContains: contains any', MustContains, [['@', '!', '$'], 2], 'thomas', {"keys": {"atleast": 2, "chars": ["@", "!", "$"]}, "msg": 'words must contains some characters'}],
    ['Test MustContains: contains only 1', MustContains, [['@', '!', '$'], 2], 'thomas@', {"keys": {"atleast": 2, "chars": ["@", "!", "$"]}, "msg": 'words must contains some characters'}],
    ['Test MustContains: valid', MustContains, [['@', '!', '$'], 2], 'thomas@$', null]
];


describe('Test validators', () => {
    tests.map((test) => {
        it(test[0], () => {
            const v = test[1](...test[2]);
            const r = v(test[3]);
            const e = test[4];

            expect(r).toEqual(e);
        })
    })
};
