#!/usr/bin/env python3
"""
OpenAPI Specification Validator for Example Go Server
Validates compliance with Constitutional Principle I: Contract-First Development

Usage: python3 .specify/scripts/validate-openapi.py
"""

import yaml
import re
import sys
from pathlib import Path

def validate_openapi_spec(spec_path):
    """Validate OpenAPI specification for constitutional compliance."""
    
    print("🔍 Validating OpenAPI Specification...")
    print(f"📄 File: {spec_path}\n")
    
    errors = []
    warnings = []
    
    # Read the file
    try:
        with open(spec_path, 'r') as f:
            content = f.read()
            spec = yaml.safe_load(content)
    except yaml.YAMLError as e:
        print(f"❌ YAML Parsing Error: {e}")
        return False
    except FileNotFoundError:
        print(f"❌ File not found: {spec_path}")
        return False
    
    print("✅ YAML syntax is valid\n")
    
    # Check 1: operationID (incorrect) vs operationId (correct)
    print("🔍 Checking operationId casing...")
    incorrect_operation_id = re.findall(r'operationID:', content)
    if incorrect_operation_id:
        errors.append(f"Found {len(incorrect_operation_id)} instances of incorrect 'operationID:' (should be 'operationId:')")
    else:
        print("✅ All operationId properties use correct lowercase 'Id'\n")
    
    # Check 2: Parameters (incorrect) vs parameters (correct)
    print("🔍 Checking parameters casing...")
    incorrect_params = re.findall(r'\n  Parameters:', content)
    if incorrect_params:
        errors.append(f"Found {len(incorrect_params)} instances of incorrect 'Parameters:' (should be 'parameters:')")
    else:
        print("✅ All parameters properties use correct lowercase\n")
    
    # Check 3: All paths have operationId
    print("🔍 Checking all endpoints have operationId...")
    paths_without_op = []
    total_endpoints = 0
    for path, methods in spec['paths'].items():
        for method, details in methods.items():
            if method != 'parameters' and isinstance(details, dict):
                total_endpoints += 1
                if 'operationId' not in details:
                    paths_without_op.append(f"{method.upper()} {path}")
    
    if paths_without_op:
        errors.append(f"Found {len(paths_without_op)} endpoints without operationId: {', '.join(paths_without_op)}")
    else:
        print(f"✅ All {total_endpoints} endpoints have operationId defined\n")
    
    # Check 4: User endpoints have userId parameter
    print("🔍 Checking user endpoints have userId parameter...")
    user_endpoints = [
        ('get', '/user/{userId}'),
        ('delete', '/user/{userId}'),
        ('get', '/user/{userId}/points'),
        ('post', '/user/{userId}/points')
    ]
    
    missing_params = []
    for method, path in user_endpoints:
        endpoint_details = spec['paths'].get(path, {}).get(method, {})
        params = endpoint_details.get('parameters', [])
        has_userId = any(p.get('name') == 'userId' for p in params if isinstance(p, dict))
        if not has_userId:
            missing_params.append(f"{method.upper()} {path}")
    
    if missing_params:
        errors.append(f"Found {len(missing_params)} user endpoints missing userId parameter: {', '.join(missing_params)}")
    else:
        print(f"✅ All {len(user_endpoints)} user endpoints have userId parameter documented\n")
    
    # Check 5: Response schemas have required fields where appropriate
    print("🔍 Checking response schemas...")
    responses_to_check = [
        ('GET', '/products', '200', ['products', 'total', 'limit']),
        ('GET', '/orders', '200', ['orders', 'total']),
        ('GET', '/user/{userId}/points', '200', ['loyaltyPoints']),
        ('POST', '/user/{userId}/points', '200', ['remainingPoints']),
    ]
    
    schemas_missing_required = []
    for method, path, status, expected_required in responses_to_check:
        endpoint = spec['paths'].get(path, {}).get(method.lower(), {})
        response = endpoint.get('responses', {}).get(status, {})
        schema = response.get('content', {}).get('application/json', {}).get('schema', {})
        required = schema.get('required', [])
        
        missing = [field for field in expected_required if field not in required]
        if missing:
            schemas_missing_required.append(f"{method} {path} {status}: missing required fields {missing}")
    
    if schemas_missing_required:
        warnings.append(f"Found {len(schemas_missing_required)} response schemas with potentially missing required fields")
        for schema in schemas_missing_required:
            print(f"  ⚠️  {schema}")
    else:
        print(f"✅ All checked response schemas have required fields defined\n")
    
    # Print summary
    print("\n" + "="*70)
    print("VALIDATION SUMMARY")
    print("="*70)
    
    if errors:
        print(f"\n❌ FAILED: {len(errors)} error(s) found:")
        for i, error in enumerate(errors, 1):
            print(f"  {i}. {error}")
    
    if warnings:
        print(f"\n⚠️  {len(warnings)} warning(s):")
        for i, warning in enumerate(warnings, 1):
            print(f"  {i}. {warning}")
    
    if not errors and not warnings:
        print("\n✅ SUCCESS: OpenAPI specification is fully compliant!")
        print(f"   - {total_endpoints} endpoints validated")
        print(f"   - OpenAPI version: {spec['openapi']}")
        print(f"   - API version: {spec['info']['version']}")
        return True
    elif not errors:
        print("\n✅ PASSED: OpenAPI specification validation passed with warnings")
        return True
    else:
        print("\n❌ FAILED: OpenAPI specification has errors that must be fixed")
        return False

if __name__ == "__main__":
    # Determine spec path relative to script location
    script_dir = Path(__file__).parent
    repo_root = script_dir.parent.parent
    spec_path = repo_root / "api" / "openapi.yaml"
    
    success = validate_openapi_spec(spec_path)
    sys.exit(0 if success else 1)
