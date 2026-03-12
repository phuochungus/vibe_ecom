.PHONY: infra-up infra-down be-new-infra-up be-new-infra-down be-mono-infra-up be-mono-infra-down

infra-up:
	$(MAKE) -C BE_new infra-up

infra-down:
	$(MAKE) -C BE_new infra-down

be-new-infra-up:
	$(MAKE) -C BE_new infra-up

be-new-infra-down:
	$(MAKE) -C BE_new infra-down

be-mono-infra-up:
	$(MAKE) -C BE_mono infra-up

be-mono-infra-down:
	$(MAKE) -C BE_mono infra-down