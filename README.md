# NpmNix

> Checkout my blog post [https://fzakaria.com/2025/03/11/nix-dynamic-derivations-a-practical-application](https://fzakaria.com/2025/03/11/nix-dynamic-derivations-a-practical-application) for an overview on this example.

This is a demonstratin of the power of _dynamic-derivations_ in Nix leveraged to build a Node project via Makefile.

⚠️ As of _2025-03-11_, you need to use [nix@d904921](https://github.com/NixOS/nix/commit/d904921eecbc17662fef67e8162bd3c7d1a54ce0) in order to use _dynamic-derivations_. Additionally, you need to enable `experimental-features = ["nix-command" "dynamic-derivations" "ca-derivations" "recursive-nix"]`. Here, there be dragons 🐲.

```console
# let's do everything in /tmp/dyn-drvs as a temporary
# nix store.
# 
# enter a bind mount for the temporary store
> nix run nixpkgs#fish --store /tmp/dyn-drvs

> nix build -f default.nix --store /tmp/dyn-drvs --print-out-paths -L
/nix/store/x9l8m94a2g6zkszab11na5l7c18xv0j1-node_modules 

> ./result
is-number  is-odd  left-pad
```

We can validate this is a correct `node_modules` by running `yarn`.

```console
> ln -s /nix/store/x9l8m94a2g6zkszab11na5l7c18xv0j1-node_modules node_modules

> yarn check --verify-tree
yarn check v1.22.22
warning package.json: No license field
success Folder in sync.
Done in 0.03s.
```

## Testing

You can run the parser by itself to generate the Nix expression that generates all the object files.

Afterwards you can build the nix expression to validate that it works.

```sh
# generate the .d files we need
> make deps
# generate the nix expression
> go run parser/parser.go ./package-lock.json > test.nix
# test it!
> nix build -f test.nix --arg pkgs 'import <nixpkgs> {}'
```